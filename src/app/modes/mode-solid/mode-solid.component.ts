import { Component, OnInit } from '@angular/core';
import { ModeSolidService } from './mode-solid.service'

@Component({
  selector: 'app-mode-solid',
  templateUrl: './mode-solid.component.html',
  styleUrls: ['./mode-solid.component.scss', '../../app.component.scss']
})
export class ModeSolidComponent implements OnInit {
  constructor(public service: ModeSolidService) { }
  ngOnInit(): void { }
}
