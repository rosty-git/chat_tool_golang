import {
  HttpEvent,
  HttpInterceptorFn,
  HttpRequest,
} from '@angular/common/http';
import { of, tap } from 'rxjs';

const TTL = 10_000;

interface CacheEntry {
  url: string;
  value: HttpEvent<unknown>;
  expiresOn: number;
}

let cacheArr = new Array<CacheEntry>();

const isCacheable = (req: HttpRequest<unknown>) => req.method === 'GET';

export const cachingInterceptor: HttpInterceptorFn = (req, next) => {
  console.log(req);
  const { urlWithParams } = req;

  if (!isCacheable(req)) {
    return next(req);
  }

  const cached = cacheArr.find((entry) => entry.url === urlWithParams);

  const now = new Date().getTime();

  if (cached && cached.expiresOn > now) {
    cacheArr = cacheArr.filter((i) => i.expiresOn > now);

    return of(cached.value);
  }

  return next(req).pipe(
    tap((response: HttpEvent<unknown>) => {
      cacheArr = cacheArr.filter((i) => i.url !== urlWithParams);

      cacheArr.push({
        url: urlWithParams,
        value: response,
        expiresOn: now + TTL,
      });
    }),
  );
};
