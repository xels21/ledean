import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';

import { LedsService } from './leds/leds.service';
import { ButtonService } from './button/button.service';
import { UpdateService } from './update/update.service';

import { LedDisplayComponent } from './led-display/led-display.component';
import { ButtonComponent } from './button/button.component';
import { NavigationComponent } from './navigation/navigation.component';
import { ModesComponent } from './modes/modes.component';
import { ModeSolidComponent } from './modes/mode-solid/mode-solid.component';
import { ModeSolidRainbowComponent } from './modes/mode-solid-rainbow/mode-solid-rainbow.component';
import { ModeRunningLedComponent } from './modes/mode-running-led/mode-running-led.component';
import { ModeTransitionRainbowComponent } from './modes/mode-transition-rainbow/mode-transition-rainbow.component';

@NgModule({
  imports: [
    BrowserModule,
    FormsModule,
    HttpClientModule
  ],
  declarations: [
    AppComponent,
    LedDisplayComponent,
    ButtonComponent,
    NavigationComponent,
    ModesComponent,
    ModeSolidComponent,
    ModeSolidRainbowComponent,
    ModeRunningLedComponent,
    ModeTransitionRainbowComponent
  ],
  providers: [
    LedsService,
    ButtonService,
    UpdateService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
