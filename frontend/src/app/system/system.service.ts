import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { REST_GET_SYSTEM_EXIT_URL } from '../config/const';


@Injectable({
  providedIn: 'root'
})
export class SystemService {

  constructor(private httpClient: HttpClient) { }

  public exit() {
    if (confirm("Stab it in the dark?")) {
      this.httpClient.get<void>(REST_GET_SYSTEM_EXIT_URL).subscribe(() => { })
    }
  }
}
