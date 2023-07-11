import { Component } from '@angular/core';
import { LedsService } from '../leds/leds.service';
import { ModesService } from '../modes/modes.service';


@Component({
  selector: 'app-led-display',
  templateUrl: './led-display.component.html',
  styleUrls: ['./led-display.component.scss']
})
export class LedDisplayComponent {

  constructor(public ledsService: LedsService, public modesService: ModesService) { }
}
