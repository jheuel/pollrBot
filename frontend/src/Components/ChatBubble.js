import React, { Component } from 'react';
import PropTypes from 'prop-types'
import './ChatBubble.css';
import FaceIcon from '@material-ui/icons/Face';
import PollIcon from '@material-ui/icons/Poll';
import Button from '@material-ui/core/Button';


class ChatBubble extends Component {
  getConversations(messages){
    if(messages === undefined){
      return;
    }

    const listItems = messages.map((message, index) => {
      let bubbleClass = 'left';
      let bubbleDirection = '';
      let icon = <PollIcon class='chat-icon'/>;

      if(message.type === 0){
        bubbleClass = 'right';
        bubbleDirection = "bubble-direction-reverse";
        icon = <FaceIcon class='chat-icon'/>;
      }

      let buttons = '';
      if (message.buttons) {
          buttons = <>
            {message.buttons.map((button, index) => (
                <Button
                    fullWidth
                    color="primary" disableElevation
                    variant="outlined"
                    size="small"
                    style={{marginTop: '8px', textTransform: 'none'}}
                    key={index}
                    >
                        {button}
                    </Button>
            ))}
        </>
      }
      return (
              <div className={`bubble-container ${bubbleDirection}`} key={index}>
                {icon}
                <div className={`bubble ${bubbleClass}`}>
                    {message.text}
                    {buttons}
                </div>
              </div>
          );
    });
    return listItems;
  }

  render() {
    const {props: {messages}} = this;
    const chatList = this.getConversations(messages);

    return (
      <div className="chats">
        <div className="chat-list">
          {chatList}
        </div>
      </div>
    );
  }
}

ChatBubble.propTypes = {
  messages: PropTypes.array.isRequired,
};

export default ChatBubble;