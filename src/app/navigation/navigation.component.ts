import { Component } from '@angular/core';
import { SystemService } from '../system/system.service'
import { LedsService } from '../leds/leds.service';

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {
  constructor( public systemService: SystemService, public ledsService: LedsService) { }
}
