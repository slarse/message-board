import * as React from "react";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Link from "@mui/material/Link";
import ProTip from "./ProTip";
import { Card, CardContent } from "@mui/material";

function Copyright() {
  return (
    <Typography variant="body2" color="text.secondary" align="center">
      {"Copyright © "}
      <Link color="inherit" href="https://mui.com/">
        Your Website
      </Link>{" "}
      {new Date().getFullYear()}.
    </Typography>
  );
}

export default function App() {
  return (
    <Container maxWidth="sm">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
					Welcome to this generic message board!
        </Typography>
        <Messages />
        <ProTip />
        <Copyright />
      </Box>
    </Container>
  );
}

function Messages() {
  const [messages, setMessages] = React.useState<MessageData[]>([]);

  React.useEffect(() => {
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
