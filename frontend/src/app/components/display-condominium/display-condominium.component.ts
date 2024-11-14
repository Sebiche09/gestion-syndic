import { Component, OnInit, ViewChild } from '@angular/core';
import { DisplayCondominiumService } from '../../services/display-condominium.service';
import { Router } from '@angular/router';
import { DialogModule } from 'primeng/dialog';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { CondominiumComponent } from '../condominium/condominium.component';

@Component({
    selector: 'displayCondominium',
    templateUrl: 'display-condominium.component.html',
    standalone: true,
    imports: [TableModule, ButtonModule, DialogModule, CondominiumComponent],
    providers: [DisplayCondominiumService]
})
export class DisplayCondominiumComponent implements OnInit {
    displayDialog: boolean = false;

    constructor(
        private displayCondominiumService: DisplayCondominiumService, 
        private router: Router
    ) {}

    ngOnInit(): void {
    }

    openDialog(): void {
        this.displayDialog = true;
    }

    closeDialog(): void {
        this.displayDialog = false;
    }
}
