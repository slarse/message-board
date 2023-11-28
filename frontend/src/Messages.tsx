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
}

export function Messages(props: MessagesProps) {
  const [messages, setMessages] = useState<MessageData[]>([]);

	const { actions } = props;

  useEffect(() => {
		actions.loadRootMessages().then(setMessages);
  }, []);

  const rootMessages = messages.filter(
    (message) => message.parentId === undefined,
  );
  const comments = buildCommentsByMessageId(messages);

  return (
    <MessageTree rootMessages={rootMessages} comments={comments} level={1} />
  );
}

type CommentsByMessageId = Map<number, MessageData[]>;

function buildCommentsByMessageId(messages: MessageData[]): CommentsByMessageId {
  const commentsByMessageId = new Map();
  for (const message of messages) {
    if (message.parentId === undefined) {
      continue;
    }

    const comments = commentsByMessageId.get(message.parentId) ?? [];
    comments.push(message);
    commentsByMessageId.set(message.parentId, comments);
  }

  return commentsByMessageId;
}

type MessageTreeProps = {
  rootMessages: MessageData[];
  comments: CommentsByMessageId;
  level: number;
};

function MessageTree(props: MessageTreeProps) {
  const { rootMessages, comments, level } = props;

  return rootMessages.length == 0 ? null : (
    <Box>
      {rootMessages.map((message) => (
        <Card sx={{ marginBottom: 2, marginTop: 2 }}>
          <MessageOverview message={message} />
          <CardActions sx={{ margin: 2 }}>
            <Button startIcon={<AddComment />}>Reply</Button>
            <Button startIcon={<AddCircle />}>Load Comments</Button>
            <Button startIcon={<Delete />}>Delete</Button>
          </CardActions>
          <Card sx={{ paddingLeft: 2 * level }}>
            <MessageTree
              rootMessages={comments.get(message.id) || []}
              comments={comments}
              level={level + 1}
            />
          </Card>
        </Card>
      ))}
    </Box>
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
