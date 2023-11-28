export type MessageData = {
  id: number;
  parentId?: number;
  content: string;
  title: string;
  author: string;
  createdAt: string;
};

export type CreateMessageData = {
  content: string;
  title: string;
  author: string;
};

export type CreateCommentData = Omit<CreateMessageData, "title"> & {
  parentId: number;
};

export type MessageActions = {
  loadMessages: () => Promise<MessageData[]>;
  loadComments: (messageId: number) => Promise<MessageData[]>;
  reply: (data: CreateCommentData) => Promise<MessageData>;
  newPost: (data: CreateMessageData) => Promise<MessageData>;
  delete: (messageId: number) => Promise<MessageData>;
};

const loadRootMessages = async () => {
  const response = await fetch("/api/messages");
  const data = await response.json();
  return data;
};

const loadComments = async (messageId: number) => {
  const response = await fetch(`/api/messages/${messageId}/comments`);
  const data = await response.json();
  return data;
};

const createMessage = async (data: CreateCommentData | CreateMessageData) => {
  const response = await fetch(`/api/messages`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data),
  });
  const responseData = await response.json();
  return responseData;
};

const deleteMessage = async (messageId: number) => {
  const response = await fetch(`/api/messages/${messageId}`, {
    method: "DELETE",
  });
  const responseData = await response.json();
  return responseData;
};

export const defaultMessageActions: MessageActions = {
  loadMessages: loadRootMessages,
  loadComments: loadComments,
  reply: createMessage,
  newPost: createMessage,
  delete: deleteMessage,
};
