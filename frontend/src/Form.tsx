import { useState } from "react";
import { Box, Button, TextField } from "@mui/material";

type CommentFormProps = {
  apply: (author: string, content: string) => Promise<void>;
};

export function CommentForm(props: CommentFormProps) {
  const [content, setContent] = useState("");

  return (
    <Box
      component="form"
      noValidate
      autoComplete="off"
      sx={{
        display: "flex",
        flexDirection: "column",
        p: 2,
        m: 2,
      }}
    >
      <TextField
        id="content"
        label="Add a reply..."
        multiline
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={4}
        sx={{ m: 1 }}
      />
      <Button
        variant="contained"
        onClick={() => props.apply(getAuthor(), content)}
        sx={{ m: 1 }}
      >
        Submit
      </Button>
    </Box>
  );
}

type MessageFormProps = {
  apply: (author: string, title: string, content: string) => Promise<void>;
};

export function MessageForm(props: MessageFormProps) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  return (
    <Box
      component="form"
      noValidate
      autoComplete="off"
      sx={{
        display: "flex",
        flexDirection: "column",
        p: 2,
      }}
    >
      <TextField
        id="title"
        label="Title"
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        sx={{ m: 1 }}
      />
      <TextField
        id="content"
        label="Message"
        multiline
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={8}
        sx={{ m: 1 }}
      />
      <Button
        variant="contained"
        onClick={() => props.apply(getAuthor(), title, content)}
        sx={{ m: 1 }}
      >
        Submit
      </Button>
    </Box>
  );
}

function getAuthor() {
  const searchParams = new URLSearchParams(window.location.search);
  return searchParams.get("author") || "John";
}
