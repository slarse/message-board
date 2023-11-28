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
    actions.loadRootMessages().then(setMessages);
  }, []);

  return (
    <Box>
      {messages.map((message) => (
        <MessageTree
          key={message.id}
          message={message}
          actions={actions}
          level={1}
        />
      ))}
    </Box>
  );
}

type MessageTreeProps = {
  message: MessageData;
  level: number;
  actions: MessageActions;
};

function MessageTree(props: MessageTreeProps) {
  const { message, level, actions } = props;
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
        <Button startIcon={<Delete />}>Delete</Button>
      </CardActions>
      <Card sx={{ paddingLeft: 2 * level, paddingRight: 2 * level }}>
        {comments.map((message) => (
          <MessageTree
            key={message.id}
            message={message}
            level={level + 1}
            actions={actions}
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
