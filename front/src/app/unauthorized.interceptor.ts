import { HttpErrorResponse, HttpInterceptorFn } from '@angular/common/http';
import { catchError, throwError } from 'rxjs';

export const unauthorizedInterceptor: HttpInterceptorFn = (req, next) => next(req).pipe(
  catchError((err: unknown) => {
    if (err instanceof HttpErrorResponse) {
      // Handle HTTP errors
      if (err.status === 401) {
        // Specific handling for unauthorized errors
        console.error('Unauthorized request:', err);

        window.location.href = '/login';
      } else {
        // Handle other HTTP error codes
        console.error('HTTP error:', err);
      }
    } else {
      // Handle non-HTTP errors
      console.error('An error occurred:', err);
    }

    // Re-throw the error to propagate it further
    return throwError(() => err);
  }),
);
