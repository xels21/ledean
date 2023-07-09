import { Component, OnInit } from '@angular/core';
import { RGB } from '../../color/color';
import { REST_MODE_SOLID_URL } from '../../config/const';
import { UpdateService } from '../../update/update.service';
import { HttpClient } from '@angular/common/http';
import { deepCopy } from '../../lib/deep-copy';
import { deepEqual } from 'fast-equals';

import { ModeSolidParameter, ModeSolidLimits, ModeSolidService } from './mode-solid.service'


@Component({
  selector: 'app-mode-solid',
  templateUrl: './mode-solid.component.html',
  styleUrls: ['./mode-solid.component.scss', '../../app.component.scss']
})
export class ModeSolidComponent implements OnInit {
  constructor(public service: ModeSolidService) {
  }

  ngOnInit(): void {}

}
