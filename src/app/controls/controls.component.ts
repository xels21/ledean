import { Component, OnInit } from '@angular/core';
import { ButtonService } from '../button/button.service';


@Component({
  selector: 'app-controls',
  templateUrl: './controls.component.html',
  styleUrls: ['./controls.component.scss']
})
export class ControlsComponent implements OnInit {

  constructor(public buttonService: ButtonService) { }

  ngOnInit(): void {
  }

}
