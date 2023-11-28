import { AddCircle, AddComment, Delete } from "@mui/icons-material";
import {
  Avatar,
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
import { CommentForm } from "./Form";
import { red } from "@mui/material/colors";

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
  const [showCommentForm, setShowCommentForm] = useState(false);

  const loadComments = async () => {
    const comments = await actions.loadComments(message.id);
    setComments(comments);
  };

  const reply = async (content: string) => {
    const searchParams = new URLSearchParams(window.location.search);
    const author = searchParams.get("author") || "John";

    setShowCommentForm(false);
    await actions.reply({
      parentId: message.id,
      content,
      author,
    });
    await loadComments();
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
        <Button
          startIcon={<AddComment />}
          onClick={() => setShowCommentForm((show) => !show)}
        >
          Reply
        </Button>
        <Button startIcon={<AddCircle />} onClick={loadComments}>
          Load Comments
        </Button>
        <Button startIcon={<Delete />} onClick={onDelete}>
          Delete
        </Button>
      </CardActions>
      {showCommentForm && <CommentForm apply={reply} />}
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
	const date = message.createdAt.substring(0, 10).replaceAll("T", " ");

  return (
    <>
      <CardHeader
        avatar={
          <Avatar sx={{ bgcolor: red[500] }} aria-label="message">
            {message.author[0]}
          </Avatar>
        }
        title={message.title}
        subheader={`${message.author} - ${date}`}
      />
      <CardContent>
        <Typography variant="body1">{message.content}</Typography>
      </CardContent>
    </>
  );
}
