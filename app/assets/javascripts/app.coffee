#= require javascripts/jquery-2.1.3.min
#= require javascripts/jquery.turbolinks.min
#= require javascripts/jquery-ujs
#= require javascripts/bootstrap.min
#= require javascripts/turbolinks
#= require javascripts/underscore
#= require javascripts/backbone
window.App =
  # Use this method to redirect so that it can be stubbed in test
  gotoUrl: (url) ->
    # Turbolinks.visit(url)
    location.href = url

  initDropdown : () ->
    $("body").on 'click', '.md-dropdown .dropdown-menu li', (event) ->
      $target = $(event.currentTarget)
      $target.closest('.input-group-btn')
             .find('[data-bind="value"]')
             .val($target.data("id")).end()
      .find('[data-bind="label"]')
      .text($target.text()).end()
      .children( '.dropdown-toggle' ).dropdown( 'toggle' )
      return false

  initWebSocket: (user_id, token) ->
    gotalk.handleNotification 'notify', (params) ->
      console.log "received #{params.count}"
    s = gotalk.connection().on 'open', ->
      console.log "gotalk connected."
    .on 'close', ->
      console.log "gotalk disconnected."

TopicDetailView = Backbone.View.extend
  el: ".topic-detail"

  events:
    "click .panel-footer a.watch": "toggleWatch"
    "click .panel-footer a.star": "toggleStar"

  toggleStar: (e) ->
    a = $(e.target)
    topicId = a.data("id")
    labelText = "#{a.data("count")} 人收藏"
    if a.hasClass("followed")
      $.post("/topics/#{topicId}/unstar").done (res) ->
        a.removeClass("followed")
        .html('<i class="fa fa-star-o"></i> ' + labelText)
    else
      $.post("/topics/#{topicId}/star").done (res) ->
        a.addClass("followed")
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

window.Topics =
  repliesPerPage: 50

  # Given floor, calculate which page this floor is in
  pageOfFloor: (floor) ->
    Math.floor((floor - 1) / Topics.repliesPerPage) + 1

  # 跳到指定楼。如果楼层在当前页，高亮该层，否则跳转到楼层所在页面并添
  # 加楼层的 anchor。返回楼层 DOM Element 的 jQuery 对象
  #
  # -   floor: 回复的楼层数，从1开始
  gotoFloor: (floor) ->
    replyEl = $("#reply#{floor}")

    if replyEl.length > 0
      Topics.highlightReply(replyEl)
    else
      page = Topics.pageOfFloor(floor)
      # TODO: merge existing query string
      url = window.location.pathname + "?page=#{page}" + "#reply#{floor}"
      App.gotoUrl url

    replyEl

  # 高亮指定楼。取消其它楼的高亮
  #
  # -   replyEl: 需要高亮的 DOM Element，须要 jQuery 对象
  highlightReply: (replyEl) ->
    $("#replies .reply").removeClass("light")
    replyEl.addClass("light")

  # 回复
  reply : (floor, login) ->
    reply_body = $(".reply-form textarea")
    new_text = "##{floor}楼 @#{login} "
    if reply_body.val().trim().length == 0
      new_text += ''
    else
      new_text = "\n#{new_text}"
    reply_body.focus().val(reply_body.val() + new_text)
    return false

  init: ->
    new TopicDetailView()

    $("#replies").on 'click', "a.mention-floor", (e) ->
      $el = $(e.target)
      floor = $el.data('floor')
      Topics.gotoFloor(floor)

    $("#replies").on "click", ".reply .btn-reply", (e) ->
      $el = $(e.target)
      Topics.reply($el.data("floor"), $el.data("login"))
    	return false



$(document).on "ready page:load", ->
  App.initDropdown()
  # App.initWebSocket()
  Topics.init()