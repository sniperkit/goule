window.goule = {} if not window.goule?

class HeaderControls
  constructor: ->
    @tabs = []
    @selected = null
  
  show: (animate) ->
    if animate
      $('#header-controls').fadeIn()
    else
      $('#header-controls').css display: 'block', opacity: '1.0'
  
  hide: (animate) ->
    if animate
      $('#header-controls').fadeOut()
    else
      $('#header-controls').css display: 'none'
  
  selectTab: (idx, animate, animateShow = true) ->
    return if @selected is idx
    @selected = idx
    tab = @tabs[idx]
    width = tab.outerWidth() + 10
    right = $(window).outerWidth() - (width + tab.offset().left) - 5
    if animate
      attributes = 'right': right, 'width': width
      $('#tab-selector').animate attributes, 'fast', ->
    else
      $('#tab-selector').css 'width': width + 'px', 'right': right + 'px'
    if idx is 0
      window.goule.services.show animateShow
      window.goule.adminSettings.hide animateShow
    else
      window.goule.adminSettings.show animateShow
      window.goule.services.hide animateShow

window.goule.headerControls = new HeaderControls()

$ ->
  $('.generate').mouseenter (e) ->
    number = Math.floor (Math.random() * (8192 - 1024)) + 1024
    $('.generate .random').html '' + number
  tabServices = $ '#tab-services'
  tabSettings = $ '#tab-settings'
  window.goule.headerControls.tabs = [tabServices, tabSettings]
  tabServices.click -> window.goule.headerControls.selectTab 0, true
  tabSettings.click -> window.goule.headerControls.selectTab 1, true