import { Component } from '@angular/core';
import { SystemService } from '../system/system.service'

@Component({
  selector: 'app-navigation',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss']
})
export class NavigationComponent {
  constructor( public systemService: SystemService) { }
}
