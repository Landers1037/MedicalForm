export namespace main {
	
	export class Doctor {
	    id: number;
	    name: string;
	    sex: string;
	    level: string;
	    phone: string;
	    email: string;
	    remark: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Doctor(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sex = source["sex"];
	        this.level = source["level"];
	        this.phone = source["phone"];
	        this.email = source["email"];
	        this.remark = source["remark"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Medicine {
	    id: number;
	    name: string;
	    manufacturer: string;
	    price: string;
	    specification: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Medicine(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.manufacturer = source["manufacturer"];
	        this.price = source["price"];
	        this.specification = source["specification"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Patient {
	    id: number;
	    date: string;
	    time: string;
	    name: string;
	    sex: string;
	    age: string;
	    ill_time: string;
	    phone: string;
	    contact: string;
	    parent: string;
	    work: string;
	    address: string;
	    allergy_history: string;
	    detail: string;
	    solution: string;
	    medical_advice: string;
	    addon: string;
	    money: string;
	    doc: string;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new Patient(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = source["date"];
	        this.time = source["time"];
	        this.name = source["name"];
	        this.sex = source["sex"];
	        this.age = source["age"];
	        this.ill_time = source["ill_time"];
	        this.phone = source["phone"];
	        this.contact = source["contact"];
	        this.parent = source["parent"];
	        this.work = source["work"];
	        this.address = source["address"];
	        this.allergy_history = source["allergy_history"];
	        this.detail = source["detail"];
	        this.solution = source["solution"];
	        this.medical_advice = source["medical_advice"];
	        this.addon = source["addon"];
	        this.money = source["money"];
	        this.doc = source["doc"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Template {
	    id: number;
	    name: string;
	    type: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new Template(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.content = source["content"];
	    }
	}

}

