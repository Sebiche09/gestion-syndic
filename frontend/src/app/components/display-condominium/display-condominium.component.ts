import { Component, OnInit, ViewChild } from '@angular/core';
import { Customer, Representative } from '../../domain/diplayCondominium.domain';
import { DisplayCondominiumService } from '../../services/display-condominium.service';
import { DataViewModule } from 'primeng/dataview';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { CommonModule } from '@angular/common';
import { TableModule } from 'primeng/table';
import { Table } from 'primeng/table';
import { HttpClientModule } from '@angular/common/http';
import { InputTextModule } from 'primeng/inputtext';
import { ProgressBarModule } from 'primeng/progressbar';
import { MultiSelectModule } from 'primeng/multiselect';
import { DropdownModule } from 'primeng/dropdown';
import { NgModel, FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { DialogModule } from 'primeng/dialog';
import { CondominiumComponent } from '../condominium/condominium.component';

@Component({
    selector: 'displayCondominium',
    templateUrl: 'display-condominium.component.html',
    standalone: true,
    imports: [TableModule, HttpClientModule, CommonModule, InputTextModule, TagModule, 
    DropdownModule, MultiSelectModule, ProgressBarModule, ButtonModule, FormsModule, DialogModule, CondominiumComponent ],
    providers: [DisplayCondominiumService, NgModel]
})
export class DisplayCondominiumComponent implements OnInit {
    displayDialog: boolean = false;
    @ViewChild('condo') condoComponent!: CondominiumComponent;

    customers!: Customer[];
    representatives!: Representative[];

    statuses!: any[];

    loading: boolean = true;

    activityValues: number[] = [0, 100];

    searchValue: string | undefined;

    constructor(private DisplayCondominiumService: DisplayCondominiumService, private router : Router) {}

    ngOnInit() {
        this.DisplayCondominiumService.getCustomersLarge().then((customers) => {
            this.customers = customers;
            this.loading = false;

            this.customers.forEach((customer) => (customer.date = new Date(<Date>customer.date)));
        });

        this.representatives = [
            { name: 'Amy Elsner', image: 'amyelsner.png' },
            { name: 'Anna Fali', image: 'annafali.png' },
            { name: 'Asiya Javayant', image: 'asiyajavayant.png' },
            { name: 'Bernardo Dominic', image: 'bernardodominic.png' },
            { name: 'Elwin Sharvill', image: 'elwinsharvill.png' },
            { name: 'Ioni Bowcher', image: 'ionibowcher.png' },
            { name: 'Ivan Magalhaes', image: 'ivanmagalhaes.png' },
            { name: 'Onyama Limba', image: 'onyamalimba.png' },
            { name: 'Stephen Shaw', image: 'stephenshaw.png' },
            { name: 'Xuxue Feng', image: 'xuxuefeng.png' }
        ];

        this.statuses = [
            { label: 'Unqualified', value: 'unqualified' },
            { label: 'Qualified', value: 'qualified' },
            { label: 'New', value: 'new' },
            { label: 'Negotiation', value: 'negotiation' },
            { label: 'Renewal', value: 'renewal' },
            { label: 'Proposal', value: 'proposal' }
        ];
    }

    clear(table: Table) {
        table.clear();
        this.searchValue = ''
    }
    openDialog() {
        this.displayDialog = true;
    }
    closeDialog() {
        this.displayDialog = false;
      }
    handleDialogHide() {
        // Appeler la méthode resetActiveIndex lorsque le dialog se ferme
        this.condoComponent.resetActiveIndex();
      }
    getSeverity(status: string) {
        switch (status.toLowerCase()) {
            case 'unqualified':
                return 'danger';

            case 'qualified':
                return 'success';

            case 'new':
                return 'info';

            case 'negotiation':
                return 'warning';

            case 'renewal':
                return null;
            default:
              return null;
        }
    }
}