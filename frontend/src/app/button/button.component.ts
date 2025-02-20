import { Component, OnInit } from '@angular/core';
import { ButtonService } from '../button/button.service';


@Component({
  selector: 'app-controls',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.scss']
})
export class ButtonComponent implements OnInit {

  constructor(public buttonService: ButtonService) { }

  ngOnInit(): void {
  }

}
