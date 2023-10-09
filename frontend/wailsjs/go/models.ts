export namespace main {
	
	export class BoolReturn {
	    success: boolean;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new BoolReturn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	    }
	}

}

