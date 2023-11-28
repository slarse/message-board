import { useState } from "react";
import { Box, Button, TextField } from "@mui/material";

type CommentFormProps = {
  apply: (content: string) => Promise<void>;
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
      />
      <Button variant="contained" onClick={() => props.apply(content)}>Submit</Button>
    </Box>
  );
}

type MessageFormProps = {
	apply: (title: string, content: string) => Promise<void>;
	title: string;
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
			/>
      <TextField
        id="content"
        label="Message"
        multiline
        value={content}
        onChange={(e) => setContent(e.target.value)}
        rows={8}
      />
      <Button variant="contained" onClick={() => props.apply(title, content)}>Submit</Button>
    </Box>
  );
}
