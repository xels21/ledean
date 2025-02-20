import { Component, OnInit } from '@angular/core';
import { ModeGradientService } from './mode-gradient.service';

@Component({
  selector: 'app-mode-gradient',
  templateUrl: './mode-gradient.component.html',
  styleUrls: ['./mode-gradient.component.scss']
})
export class ModeGradientComponent implements OnInit {
  constructor(public service: ModeGradientService) { }
  ngOnInit(): void { }
}
