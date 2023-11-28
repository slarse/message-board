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

export type CreateCommentData = CreateMessageData & {
  parentId: number;
};

export type MessageActions = {
  loadRootMessages: () => Promise<MessageData[]>;
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

export const defaultMessageActions: MessageActions = {
  loadRootMessages: loadRootMessages,
  loadComments: loadComments,
  reply: (data: CreateCommentData) => Promise.resolve({} as MessageData),
  newPost: (data: CreateMessageData) => Promise.resolve({} as MessageData),
  delete: (messageId: number) => Promise.resolve({} as MessageData),
};
