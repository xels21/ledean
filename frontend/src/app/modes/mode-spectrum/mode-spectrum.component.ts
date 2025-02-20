import { Component, OnInit } from '@angular/core';
import { ModeSpectrumService } from './mode-spectrum.service';

@Component({
  selector: 'app-mode-spectrum',
  templateUrl: './mode-spectrum.component.html',
  styleUrls: ['./mode-spectrum.component.scss', '../../app.component.scss']
})
export class ModeSpectrumComponent implements OnInit {
  constructor(public service: ModeSpectrumService) { }
  ngOnInit(): void { }
}
