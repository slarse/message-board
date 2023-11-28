import { AddCircle, AddComment, Delete } from "@mui/icons-material";
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardHeader,
  Typography,
} from "@mui/material";
import { useEffect, useState } from "react";
import { MessageActions, MessageData } from "./MessageActions";

type MessagesProps = {
  actions: MessageActions;
};

export function Messages(props: MessagesProps) {
  const [messages, setMessages] = useState<MessageData[]>([]);

  const { actions } = props;

  useEffect(() => {
    actions.loadMessages().then(setMessages);
  }, []);

  const onDelete = async (messageId: number) => {
    const redactedMessage = await actions.delete(messageId);
    setMessages(
      messages.map((message) =>
        message.id === messageId ? redactedMessage : message,
      ),
    );
  };

  return (
    <Box>
      {messages.map((message) => (
        <MessageTree
          key={message.id}
          message={message}
          actions={actions}
          level={1}
          onDelete={() => onDelete(message.id)}
        />
      ))}
    </Box>
  );
}

type MessageTreeProps = {
  message: MessageData;
  level: number;
  actions: MessageActions;
  onDelete: () => Promise<void>;
};

function MessageTree(props: MessageTreeProps) {
  const { message, level, actions, onDelete } = props;
  const [comments, setComments] = useState<MessageData[]>([]);

  const loadComments = async () => {
    const comments = await actions.loadComments(message.id);
    setComments(comments);
  };

  const reply = async () => {
    const comment = await actions.reply({
      parentId: message.id,
      content: `Re: ${message.content}`,
      author: "Paul",
    });
    setComments([...comments, comment]);
  };

  const onDeleteComment = async (messageId: number) => {
    const redactedComment = await actions.delete(messageId);
    setComments(
      comments.map((comment) =>
        comment.id === messageId ? redactedComment : comment,
      ),
    );
  };

  return (
    <Card variant="outlined" sx={{ marginBottom: 2, marginTop: 2 }}>
      <MessageOverview message={message} />
      <CardActions sx={{ margin: 2 }}>
        <Button startIcon={<AddComment />} onClick={reply}>
          Reply
        </Button>
        <Button startIcon={<AddCircle />} onClick={loadComments}>
          Load Comments
        </Button>
        <Button startIcon={<Delete />} onClick={onDelete}>
          Delete
        </Button>
      </CardActions>
      <Card sx={{ paddingLeft: 2 * level, paddingRight: 2 * level }}>
        {comments.map((message) => (
          <MessageTree
            key={message.id}
            message={message}
            level={level + 1}
            actions={actions}
            onDelete={() => onDeleteComment(message.id)}
          />
        ))}
      </Card>
    </Card>
  );
}
type MessageOverviewProps = {
  message: MessageData;
};

function MessageOverview(props: MessageOverviewProps) {
  const { message } = props;

  return (
    <>
      <CardHeader>
        <Typography variant="h4">{message.title}</Typography>
      </CardHeader>
      <CardContent>
        <Typography variant="h6">
          {message.author} {message.createdAt}
        </Typography>
        <Typography variant="body1">{message.content}</Typography>
      </CardContent>
    </>
  );
}
