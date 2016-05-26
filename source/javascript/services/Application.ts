module Application {

    export class TechnicConfig {
        public isSolderPack:number = 1;

        public createForgeZip:boolean = false;
        public forgeVersion:ForgeVersion.ForgeVersion;

        public checkPermissions:boolean = false;
        public isPublicPack:boolean = true;
    }

    export class FtbConfig {
        public isPublicPack:boolean = true;
    }

    export class Folder {
        public name:string;
        public include:boolean;
    }
    
    export class Modpack {
        public name:string;
        public inputDirectory:string = "";
        public outputDirectory:string = "";
        public clearOutputDirectory:boolean = true;
        public minecraftVersion:string = "1.9";
        public version:string = "1.0.0";
        public additionalFolders:Array<Folder> = [];
        public technic:TechnicConfig = new TechnicConfig();
        public ftb:FtbConfig = new FtbConfig();

        constructor() {
            this.name = "Unnamed modpack";
            console.log(this);
        }

        public static fromJson(data:Modpack):Modpack {
            var modpack = new Modpack();
            modpack.name = data.name;
            modpack.inputDirectory = data.inputDirectory;
            modpack.outputDirectory = data.outputDirectory;
            modpack.clearOutputDirectory = data.clearOutputDirectory;
            modpack.minecraftVersion = data.minecraftVersion;
            modpack.version = data.version;
            modpack.additionalFolders = data.additionalFolders;
            modpack.technic = data.technic;
            modpack.ftb = data.ftb;
            return modpack;
        }
        
        public isValid(): boolean {
            if(!this.name) return false;
            if(!this.inputDirectory) return false;
            if(!this.outputDirectory) return false;
            if(!this.minecraftVersion) return false;
            if(!this.version) return false;
            return true;
        }
    }

    export class Mod {
        public modid: string;
        public name: string;
        public description: string;
        public version: string;
        public mcversion: string;
        public url: string;
        public authors: string;
        public credits: string;
        public filename: string;
        // Naming is totally a hack to make sure the value does not get send to the server
        public $$isDone: boolean;

        public static fromJson(data:Mod): Mod {
            var m = new Mod();
            m.modid = data.modid;
            m.name = data.name;
            m.description = data.description;
            m.version = data.version;
            m.mcversion = data.mcversion;
            m.url = data.url;
            m.authors = data.authors;
            m.credits = data.credits;
            m.filename = data.filename;
            return m;
        }
        
        public isValid(): boolean {
            if (!this.modid) return false;
            if (!this.name) return false;
            if (!this.version) return false;
            if (!this.mcversion) return false;
            if (this.authors.length <= 0) return false;
            return true;
        }
    }


    export class Application {
        static $inject = ["$rootScope", "goComm", "$state"];
        public modpacks:Array<Modpack> = [];
        public modpack:Modpack;

        constructor(protected $rootScope:angular.IRootScopeService, protected goComm:GoCommService.GoCommService, protected $state: angular.ui.IStateService) {
            var self = this;
            goComm.send("load-modpacks", {});
            $rootScope.$watch(function () {
                return self.modpacks;
            }, function () {
                self.saveModpackData()
            }, true);
            $rootScope.$on("data-loaded", (event:angular.IAngularEvent, modpacks:Array<Modpack>) => {
                modpacks.forEach((modpack:Modpack) => {
                    self.modpacks.push(Modpack.fromJson(modpack))
                });
                if (self.modpacks.length) {
                    self.modpack = self.modpacks[0];
                    self.$state.go("modpack");
                }
            });
        }

        protected saveModpackData():void {
            this.goComm.send("save-modpacks", this.modpacks);
        }
    }


    angular.module("ModpackHelper").service("application", Application);
}
