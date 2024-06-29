import { Component, ElementRef, ViewChild } from '@angular/core';
import * as pdfjsLib from 'pdfjs-dist';

import { environment } from '../../environments/environment';
import { DataService, GalleryState } from '../data.service';

@Component({
  selector: 'app-pdf-viewer',
  standalone: true,
  imports: [],
  templateUrl: './pdf-viewer.component.html',
  styleUrls: ['./pdf-viewer.component.scss'],
})
export class PdfViewerComponent {
  @ViewChild('pdfViewer', { static: false })
    pdfViewer!: ElementRef<HTMLDivElement>;

  state: GalleryState = {
    isOpen: false,
  };

  s3Prefix = environment.S3_PREFIX;

  constructor(private dataService: DataService) {
    pdfjsLib.GlobalWorkerOptions.workerSrc =
      'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/4.3.136/pdf.worker.min.mjs';

    this.dataService.galleryState$.subscribe((value) => {
      this.state = value;

      console.log(value);

      if (
        value.index != null &&
        value.files?.[value.index].type === 'application/pdf'
      ) {
        console.log('pdf');

        const loadingTask = pdfjsLib.getDocument(
          this.s3Prefix + value.files[value.index].s3_key,
        );
        loadingTask.promise.then(
          (pdf) => {
            console.log(pdf);

            // eslint-disable-next-line no-plusplus
            for (let pageNum = 1; pageNum <= pdf.numPages; pageNum++) {
              // eslint-disable-next-line @typescript-eslint/no-loop-func
              pdf.getPage(pageNum).then((page) => {
                console.log(page);

                const canvas = document.createElement('canvas');
                const context = canvas.getContext('2d');
                const viewport = page.getViewport({ scale: 1.0 });

                canvas.width = viewport.width;
                canvas.height = viewport.height;

                const renderContext = {
                  canvasContext: context!,
                  viewport,
                };

                page.render(renderContext);
                this.pdfViewer.nativeElement.appendChild(canvas);
              });
            }
          },
          (reason) => console.error(reason),
        );
      }
    });
  }
}
