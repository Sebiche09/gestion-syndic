import { Component, OnInit } from '@angular/core';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { RippleModule } from 'primeng/ripple';
import { TagModule } from 'primeng/tag';
import { UnitService } from '../../services/unit.service';
@Component({
  selector: 'app-display-unit',
  standalone: true,
  imports: [TableModule, ButtonModule, RippleModule, TagModule],
  providers: [UnitService],
  templateUrl: './display-unit.component.html',
  styleUrl: './display-unit.component.scss'
})
export class DisplayUnitComponent implements OnInit{
  units: any[] = [];

  constructor(private unitService: UnitService) {}


  ngOnInit() {
    this.unitService.getUnits().subscribe((data) => {
        this.units = data;
        console.log(this.units)
    });
  }
  calculateUnitTotal(): number {
    return this.units ? this.units.length : 0;
}
}
