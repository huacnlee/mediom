class Notifier
  constructor: ->
    @enableNotification = false

  hasSupport: ->
    window.Notification?

  requestPermission: (cb) ->
    window.Notification.requestPermission (cb)

  setPermission: =>
    if @hasPermission()
      $('#notification-alert a.close').click()
      @enableNotification = true
    else if window.Notification.permission is "granted"
      $('#notification-alert a.close').click()

  hasPermission: ->
    if window.Notification.permission is "granted"
      return true
    else
      return false

  checkOrRequirePermission: =>
    if @hasSupport()
      if @hasPermission()
        @enableNotification = true
      else
        if window.Notification.permission is "default"
          @showTooltip()
    else
      console.log("Desktop notifications are not supported for this Browser/OS version yet.")

  showTooltip: ->
    console.log "show notifications tip"
    $('.main-container').prepend("<div class='alert alert-info' id='notification-alert'><a href='#' id='link_enable_notifications'>点击这里</a> 开启桌面提醒通知功能。 <a class='close fa fa-close' data-dismiss='alert' href='#'></a></div>")
    $("#notification-alert").alert()
    $('#notification-alert').on 'click', 'a#link_enable_notifications', (e) =>
      e.preventDefault()
      @requestPermission(@setPermission)

  visitUrl: (url) ->
    window.location.href = url

  notify: (avatar, title, content, url = null) ->
    if @enableNotification
      opts =
        icon: avatar
        body: content
      popup = new window.Notification(title,opts)
      popup.onclick = ->
        window.parent.focus()
        $.notifier.visitUrl(url)

jQuery.notifier = new Notifier
