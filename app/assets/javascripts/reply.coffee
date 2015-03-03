class ReplyView extends Backbone.View
  tagName: 'div'
  className: 'reply-form'
  events:
    'click .btn-primary': "submit"
    
  submit: (e) ->
    