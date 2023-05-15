export namespace main {
	
	export class AISelectResult {
	    bestMove: string;
	    bestScore: number;
	    elapsedTime: string;
	
	    static createFrom(source: any = {}) {
	        return new AISelectResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bestMove = source["bestMove"];
	        this.bestScore = source["bestScore"];
	        this.elapsedTime = source["elapsedTime"];
	    }
	}

}

