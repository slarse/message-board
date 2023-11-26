import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import { Messages } from "./Messages";

export default function App() {
  return (
    <Container maxWidth="sm">
      <Box sx={{ my: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
					Welcome to this generic message board!
        </Typography>
        <Messages />
      </Box>
    </Container>
  );
}
