import { Injectable } from '@angular/core';


export interface UpdateIntervall {
  cb();
  timeout: number;
}
@Injectable({
  providedIn: 'root'
})
export class UpdateService {

  updateArr: Array<UpdateIntervall>;
  activeIntervallArr: Array<any>;
  isActive: boolean

  constructor() {
    this.updateArr = new Array<UpdateIntervall>()
    this.activeIntervallArr = new Array<any>()
    this.isActive = true
  }

  registerPolling(updateIntervall: UpdateIntervall) {
    this.updateArr.push(updateIntervall)
    if (this.isActive) {
      this.activeIntervallArr.push(setInterval(() => updateIntervall.cb() , updateIntervall.timeout))
    }
  }

  stopPolling() {
    this.activeIntervallArr.forEach(activeIntervall => {
      clearInterval(activeIntervall);
    });
    this.activeIntervallArr = new Array<any>()
    this.isActive = false
  }

  startPolling() {
    this.updateArr.forEach(update => {
      this.activeIntervallArr.push(setInterval(() => update.cb(), update.timeout))
    });
    this.isActive = true
  }

  switchState(active) {
    if (this.isActive == active) {
      return
    }
    else if (this.isActive) {
      this.stopPolling()
    } else {
      this.startPolling()
    }
  }

  update(){
    this.updateArr.forEach(update => {
      update.cb()
    });
  }

}
