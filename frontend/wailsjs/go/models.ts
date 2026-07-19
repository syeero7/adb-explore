export namespace gadb {
	
	export class DeviceFileInfo {
	    Name: string;
	    Mode: number;
	    Size: number;
	    // Go type: time
	    LastModified: any;
	
	    static createFrom(source: any = {}) {
	        return new DeviceFileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Mode = source["Mode"];
	        this.Size = source["Size"];
	        this.LastModified = this.convertValues(source["LastModified"], null);
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

}

