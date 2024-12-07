import { Component } from '@angular/core';
import { ModePictureService } from './mode-picture.service'

@Component({
  selector: 'app-mode-picture',
  templateUrl: './mode-picture.component.html',
  styleUrls: ['./mode-picture.component.scss']
})
export class ModePictureComponent {
  constructor(public service: ModePictureService) { }
  ngOnInit(): void { }
}
