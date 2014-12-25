package goule

import (
	"encoding/json"
	"errors"
	"github.com/unixpickle/ezserver"
	"github.com/unixpickle/gohttputil"
	"github.com/unixpickle/reverseproxy"
	"net/http"
	"reflect"
)

type api struct {
	*Goule
	w http.ResponseWriter
	r *http.Request
}

// AddRuleAPI adds a new proxy rule.
func (a *api) AddRuleAPI(rule reverseproxy.Rule) {
	a.config.Rules = append(a.config.Rules, rule)
	a.config.Save()
}

// AddServiceAPI adds a new service and possibly starts it.
func (a *api) AddServiceAPI(name string, cfg Service) error {
	if _, ok := a.services[name]; ok {
		return errors.New("Service name already taken.")
	}
	
	// Create the executor.Service and possibly start it
	excService := cfg.ToExecutorService()
	a.services[name] = excService
	if cfg.Autolaunch {
		excService.Start()
	}
	
	// Update configuration
	a.config.Services[name] = cfg
	a.config.Save()
	return nil
}

// AuthAPI returns whether the given password is correct.
func (a *api) AuthAPI(password string) bool {
	if !a.config.Admin.Try(password) {
		return false
	}
	// Create a new cookie and set it.
	id := a.sessions.login()
	cookie := &http.Cookie{Name: SessionIdCookie, Value: id}
	http.SetCookie(a.w, cookie)
	return true
}

// Call performs an API.
func (a *api) Call(name string, body []byte) ([]byte, int, error) {
	// Find the method for the given API.
	method := reflect.ValueOf(a).MethodByName(name + "API")
	if !method.IsValid() {
		return nil, http.StatusNotFound, errors.New("Unknown API: " + name)
	}

	// Decode the array of JSON-encoded arguments.
	var rawArgs []string
	if err := json.Unmarshal(body, &rawArgs); err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Decode the exact arguments.
	args, err := decodeArgs(method, rawArgs)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	// Lock the mutex in the appropriate way
	if isWriteAPI(name) {
		a.mutex.Lock()
		defer a.mutex.Unlock()
	} else {
		a.mutex.RLock()
		defer a.mutex.RUnlock()
	}

	// Convert the return value to an array of serializable objects.
	res := method.Call(args)
	resList := make([]interface{}, len(res))
	for i, val := range res {
		rawValue := val.Interface()
		// Convert errors to strings
		if err, ok := rawValue.(error); ok {
			rawValue = err.Error()
		}
		resList[i] = rawValue
	}
	
	// Encode the result
	if encoded, err := json.Marshal(resList); err != nil {
		return nil, http.StatusInternalServerError, err
	} else {
		return encoded, 0, nil
	}
}

func (a *api) ConfigAPI() *Config {
	return a.config
}

// DeauthAPI does nothing.
func (a *api) DeauthAPI() {
	// Invalidate the current session
	cookie, _ := a.r.Cookie(SessionIdCookie)
	a.sessions.logout(cookie.Value)
	
	// Delete the cookie on the client-side
	content := SessionIdCookie + "=deleted; " +
		"expires=Thu, 01 Jan 1970 00:00:00 GMT"
	a.w.Header()["Set-Cookie"] = []string{content}
}

// DeleteRuleAPI deletes a rule by value
func (a *api) DeleteRuleAPI(rule reverseproxy.Rule) error {
	for i, r := range a.config.Rules {
		if rulesEqual(r, rule) {
			// Remove the rule
			a.config.Rules = append(a.config.Rules[0:i],
				a.config.Rules[i+1:]...)
			a.config.Save()
			return nil
		}
	}
	return errors.New("Rule not found.")
}

// DeleteServiceAPI deletes a service by name.
func (a *api) DeleteServiceAPI(name string) error {
	service, ok := a.services[name]
	if !ok {
		return errors.New("Service not found.")
	}
	service.Stop()
	delete(a.services, name)
	delete(a.config.Services, name)
	a.config.Save()
	return nil
}

// Handle handles the API call and writes a JSON response.
func (a *api) Handle() {
	// The path is "/api/APINAME"
	name := a.r.URL.Path[5:]
	
	// Make sure they are authorized to make this request.
	authed := a.w.Header().Get("Set-Cookie") != ""
	if !authed && name != "Auth" {
		gohttputil.RespondJSON(a.w, http.StatusForbidden, "Permissions denied.")
		return
	}

	// Read the contents of the request
	contents, err := gohttputil.ReadRequest(a.r, 0x10000)
	if err != nil {
		gohttputil.RespondJSON(a.w, http.StatusBadRequest, err.Error())
		return
	}

	// Run the call
	response, code, err := a.Call(name, contents)
	if err != nil {
		gohttputil.RespondJSON(a.w, code, err.Error())
		return
	}
	a.w.Write(response)
}

// SetAdminPortAPI updates the admin port.
func (a *api) SetAdminPortAPI(port int) error {
	a.admin.Stop()
	if err := a.admin.Start(port); err != nil {
		// Attempt to restart it on the old port.
		a.admin.Start(a.config.Admin.Port)
		return err
	}
	// Port change was successful; save configuration
	a.config.Admin.Port = port
	a.config.Save()
	return nil
}

// SetAssetsAPI sets the admin assets path.
func (a *api) SetAssetsAPI(path string) {
	a.config.Admin.Assets = path
	a.config.Save()
}

// SetPasswordAPI sets the new administrative password.
func (a *api) SetPasswordAPI(password string) {
	a.config.Admin.Hash = Hash(password)
	a.config.Save()
}

// SetRuleAPI replaces an old rule with a new rule
func (a *api) SetRuleAPI(old, rule reverseproxy.Rule) {
	for i, r := range a.config.Rules {
		if rulesEqual(r, old) {
			a.config.Rules[i] = rule
			a.config.Save()
			return
		}
	}
}

// SetSessionTimeoutAPI sets the session timeout in seconds.
func (a *api) SetSessionTimeoutAPI(timeout int) {
	a.config.Admin.Timeout = timeout
	a.config.Save()
}

// SetTLSAPI sets the TLS configuration for HTTPS.
func (a *api) SetTLSAPI(tls ezserver.TLSConfig) {
	a.https.SetTLSConfig(&tls)
	a.config.TLS = tls
	a.config.Save()
}

// StartAPI starts a service by name
func (a *api) StartAPI(name string) error {
	service, ok := a.services[name]
	if !ok {
		return errors.New("Service not found.")
	}
	return service.Start()
}

// StopAPI stops a service by name
func (a *api) StopAPI(name string) error {
	service, ok := a.services[name]
	if !ok {
		return errors.New("Service not found.")
	}
	return service.Stop()
}

// UpdateServiceAPI updates a service by name.
func (a *api) UpdateServiceAPI(name string, service Service) error {
	oldServ, ok := a.services[name]
	if !ok {
		return errors.New("Service not found.")
	}
	oldServ.Stop()
	a.services[name] = service.ToExecutorService()
	a.config.Services[name] = service
	a.config.Save()
	return nil
}

func decodeArgs(method reflect.Value, raw []string) ([]reflect.Value, error) {
	// Make sure they passed the right number of arguments
	if method.Type().NumIn() != len(raw) {
		return nil, errors.New("Invalid number of arguments.")
	}

	// Decode each argument separately.
	res := make([]reflect.Value, len(raw))
	for i, rawArg := range raw {
		inputType := method.Type().In(i)
		dec := reflect.New(inputType)
		if err := json.Unmarshal([]byte(rawArg), dec.Interface()); err != nil {
			return nil, err
		}
		res[i] = reflect.Indirect(dec)
	}

	return res, nil
}

func isWriteAPI(name string) bool {
	return name != "Config" && name != "Services"
}

func rulesEqual(r1 reverseproxy.Rule, r2 reverseproxy.Rule) bool {
	return reflect.DeepEqual(r1, r2)
}
