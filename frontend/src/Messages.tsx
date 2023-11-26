import { Box, Card, CardContent, Typography } from "@mui/material";
import { useEffect, useState } from "react";

export function Messages() {
  const [messages, setMessages] = useState<MessageData[]>([]);

  useEffect(() => {
    fetch("/api/messages")
      .then((response) => response.json())
      .then((data) => setMessages(data));
  }, []);

  return <MessageList messages={messages} />;
}

type MessageListProps = {
  messages: MessageData[];
};

function MessageList(props: MessageListProps) {
  const { messages } = props;

  return (
    <Box>
      {messages.map((message) => (
        <MessageOverview key={message.id} message={message} />
      ))}
    </Box>
  );
}

type MessageData = {
  id: number;
  content: string;
  title: string;
  author: string;
	createdAt: string;
};

type MessageOverviewProps = {
  message: MessageData;
};

function MessageOverview(props: MessageOverviewProps) {
  const { message } = props;

  return (
    <Card variant="elevation" sx={{ marginBottom: 2 }}>
      <CardContent>
        <Typography variant="h4">{message.title}</Typography>
        <Typography variant="h6">{message.author} {message.createdAt}</Typography>
        <Typography variant="body1">{message.content}</Typography>
      </CardContent>
    </Card>
  );
}
