import { Component, NgModule } from '@angular/core';
import { FloatLabelModule } from "primeng/floatlabel";
import { FormsModule } from '@angular/forms';

@Component({
    selector: 'float-label',
    templateUrl: './float-label.component.html',
    standalone: true,
    imports: [FloatLabelModule, FormsModule]
})

export class FloatLabelComponent {
    value: string | undefined;
}

