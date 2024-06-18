const BASE = 'dev-chat.studyraft.com';

export const environment = {
  production: true,
  
  BASE_API_URL: `https://${BASE}`,
  BASE_WS_URL: `wss://${BASE}`,

  POSTS_PAGE_SIZE: 20,

  S3_PREFIX: 'https://BUCKET_NAME.s3.us-east-2.amazonaws.com/',
};
