import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

import { GlobalVariable } from '../../global';
import { FrontFile } from '../data.service';

@Component({
  selector: 'app-file',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './file.component.html',
  styleUrl: './file.component.scss',
})
export class FileComponent {
  @Input() file: FrontFile = { name: '', size: 0, type: '', ext: '' };

  s3Prefix = GlobalVariable.S3_PREFIX;
}
