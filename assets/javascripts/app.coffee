#= require javascripts/jquery-2.1.3.min
#= require javascripts/jquery.turbolinks.min
#= require javascripts/jquery-ujs
#= require javascripts/bootstrap.min
#= require javascripts/turbolinks
#= require javascripts/underscore
#= require javascripts/backbone
AppView = Backbone.View.extend
  el: "document"

  repliesPerPage: 50

  events:
    "click .topic-detail .panel-footer a.watch": "toggleWatch"
    "click .topic-detail .panel-footer a.star": "toggleStar"
    "click .md-dropdown .dropdown-menu li": "toggleDropdown"
    "click #replies .reply .btn-reply": "reply"
    "click #replies a.mention-floor": "mentionFloor"
    "click .button-captcha": "refreshCaptcha"
    "keydown": "handleKeyDown"

  initialize: ->
    @initWebSocket()

  initWebSocket: ->
    @ws = new WebSocket("ws://#{window.location.host}/msg")
    @ws.onmessage = @onWebSocketMessage

  onWebSocketMessage: (res) ->
    notify = JSON.parse(res.data)
    badge = $(".notification-count .badge")
    if notify.unread_count > 0
      badge.addClass("badge-new").text(notify.unread_count)
    else
      badge.removeClass("badge-new").text(0)

  toggleDropdown: (e) ->
    $target = $(e.currentTarget)
    $target.closest('.input-group-btn')
           .find('[data-bind="value"]')
           .val($target.data("id")).end()
    .find('[data-bind="label"]')
    .text($target.text()).end()
    .children( '.dropdown-toggle' ).dropdown( 'toggle' )
    return false

  toggleStar: (e) ->
    a = $(e.target)
    topicId = a.data("id")
    count = parseInt(a.data("count"))
    if a.hasClass("followed")
      $.post("/topics/#{topicId}/unstar").done (res) ->
        newCount = count - 1
        labelText = "#{newCount} 人收藏"
        a.removeClass("followed")
        .data("count", newCount)
        .html('<i class="fa fa-star-o"></i> ' + labelText)
    else
      $.post("/topics/#{topicId}/star").done (res) ->
        newCount = count + 1
        labelText = "#{newCount} 人收藏"
        a.addClass("followed")
        .data("count", newCount)
        .html('<i class="fa fa-star"></i> ' + labelText)
    return false

  toggleWatch: (e) ->
    a = $(e.target)
    topicId = a.data("id")
    if a.hasClass("followed")
      $.post("/topics/#{topicId}/unwatch").done (res) ->
        a.removeClass("followed")
        .attr("title", "关注此话题，当有新回帖的时候会收到通知")
        .html('<i class="fa fa-eye"></i> 关注')
    else
      $.post("/topics/#{topicId}/watch").done (res) ->
        a.addClass("followed")
        .attr("title", "点击取消关注")
        .html('<i class="fa fa-eye"></i> 已关注')
    return false

  reply: (e) ->
    _el = $(e.target)
    floor = _el.data("floor")
    login = _el.data("login")
    reply_body = $(".reply-form textarea")
    new_text = "##{floor}楼 @#{login} "
    if reply_body.val().trim().length == 0
      new_text += ''
    else
      new_text = "\n#{new_text}"
    reply_body.focus().val(reply_body.val() + new_text)
    return false

  mentionFloor: (e) ->
    _el = $(e.target)
    floor = _el.data('floor')
    replyEl = $("#reply#{floor}")
    if replyEl.length > 0
      @highlightReply(replyEl)
    else
      page = @pageOfFloor(floor)
      # TODO: merge existing query string
      url = window.location.pathname + "?page=#{page}" + "#reply#{floor}"
      @gotoUrl(url)

    replyEl

  highlightReply: (replyEl) ->
    $("#replies .reply").removeClass("light")
    replyEl.addClass("light")

  pageOfFloor: (floor) ->
    Math.floor((floor - 1) / Topics.repliesPerPage) + 1

  gotoUrl: (url) ->
    # Turbolinks.visit(url)
    location.href = url

  refreshCaptcha: (e) ->
    img = $(e.target)
    img.attr("src", "/captcha?t=#{(new Date).getTime()}")
    return false

  handleKeyDown: (e) ->
    self = @
    console.log e
    if e.target.nodeName == 'BODY' or e.target.nodeName == 'HTML'
      if !e.ctrlKey && !e.altKey && !e.shiftKey && !e.metaKey
        self.pageAction(e.which)

  pageAction: (code) ->
    if code == 37
      link = $(".pager .previous a")
    elsif code == 39
      link = $(".pager .next a")
    href = link.attr("href")
    if href && href != document.location && href != "#"
      Turbolinks.visit(href)
    return false


$(document).on "ready page:load", ->
  new AppView()