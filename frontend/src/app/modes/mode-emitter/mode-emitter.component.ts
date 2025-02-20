import { Component, OnInit } from '@angular/core';
import { ModeEmitterService } from './mode-emitter.service';

@Component({
  selector: 'app-mode-emitter',
  templateUrl: './mode-emitter.component.html',
  styleUrls: ['./mode-emitter.component.scss','../../app.component.scss']
})
export class ModeEmitterComponent implements OnInit {
  constructor(public service:ModeEmitterService) { }
  ngOnInit(): void {}
}
