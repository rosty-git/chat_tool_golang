const BASE = 'localhost:8080';

export const GlobalVariable = Object.freeze({
  BASE_API_URL: `http://${BASE}`,
  BASE_WS_URL: `ws://${BASE}`,

  POSTS_PAGE_SIZE: 20,

  S3_PREFIX: 'https://cchhaatt.s3.us-east-2.amazonaws.com/',
});
