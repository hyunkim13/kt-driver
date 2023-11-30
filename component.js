define('nodes/components/driver-kt/component', ['exports', 'shared/mixins/node-driver'], function (exports, _nodeDriver) {
  'use strict';

  Object.defineProperty(exports, "__esModule", {
    value: true
  });

  var LAYOUT = "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";

  var computed = Ember.computed;
  var get = Ember.get;
  var set = Ember.set;
  var alias = Ember.computed.alias;
  var service = Ember.inject.service;

  var defaultRadix = 10;
  var defaultBase = 1024;
  exports.default = Ember.Component.extend(_nodeDriver.default, {
    driverName: 'cloudca',
    config: alias('model.cloudcaConfig'),
    app: service(),

    init: function init() {
      var decodedLayout = window.atob(LAYOUT);
      var template = Ember.HTMLBars.compile(decodedLayout, {
        moduleName: 'nodes/components/driver-kt/template'
      });
      set(this, 'layout', template);

      this._super.apply(this, arguments);
    },

    bootstrap: function bootstrap() {
      var config = get(this, 'globalStore').createRecord({
        type: 'cloudcaConfig',
        apiUrl: "https://api.cloud.ca/v1",
        usePrivateIp: true,
        sshUser: "cca-user"
      });
      set(this, 'model.cloudcaConfig', config);

      var self = this;

      set(this, "onTemplateChange", function (template) {
        self.set("config.sshUser", template.sshUser);
      });

      set(this, 'onComputeOfferingChange', function (co) {
        self.set("config.customCompute", co.custom);
        if (co.custom) {
          self.set("config.cpuCount", 1);
          self.set("config.memoryMb", 2048);
        } else {
          self.set("config.cpuCount", null);
          self.set("config.memoryMb", null);
        }
      });

      set(this, 'onDiskOfferingChange', function (o) {
        self.set("config.customDiskOffering", o.customSize);
        if (o.customSize) {
          self.set("config.additionalDiskSizeGb", 20);
        } else {
          self.set("config.additionalDiskSizeGb", null);
        }
      });
    },

    validate: function validate() {
      this._super();
      var errors = this.get('errors') || [];

      var name = this.get('model.hostname');
      if (name) {
        if (name.length > 62) {
          errors.push('Name can be a maximum of 62 characters long.');
        } else if (!/^[a-zA-Z]/.test(name) || !/[a-zA-Z0-9]$/.test(name)) {
          errors.push('Name must start with a letter and end with a letter or digit.');
        } else if (!/^[-a-zA-Z0-9]+$/.test(name)) {
          errors.push('Name can only contain letters, digits and hyphens.');
        }
      }
      this.set('errors', errors);
      return !errors.length;
    },

    apiKeyPage: true,
    pages: ['apiKeyPage', 'envPage', 'computePage', 'nodeTemplate'],

    cpuOptions: [{
      name: '1 vCPU',
      value: '1'
    }, {
      name: '2 vCPU',
      value: '2'
    }, {
      name: '4 vCPU',
      value: '4'
    }, {
      name: '6 vCPU',
      value: '6'
    }, {
      name: '8 vCPU',
      value: '8'
    }, {
      name: '10 vCPU',
      value: '10'
    }, {
      name: '12 vCPU',
      value: '12'
    }, {
      name: '16 vCPU',
      value: '16'
    }],

    memOptions: [{
      name: '2 GB',
      value: 1024 * 2 + ''
    }, {
      name: '4 GB',
      value: 1024 * 4 + ''
    }, {
      name: '8 GB',
      value: 1024 * 8 + ''
    }, {
      name: '12 GB',
      value: 1024 * 12 + ''
    }, {
      name: '16 GB',
      value: 1024 * 16 + ''
    }, {
      name: '20 GB',
      value: 1024 * 20 + ''
    }, {
      name: '24 GB',
      value: 1024 * 24 + ''
    }, {
      name: '28 GB',
      value: 1024 * 28 + ''
    }, {
      name: '32 GB',
      value: 1024 * 32 + ''
    }],

    setPage: function setPage(pageNum) {
      this.set('errors', []);
      var self = this;
      this.pages.forEach(function (p, i) {
        self.set(p, i === pageNum);
      });
    },
    actions: {
      goToApiKeyPage: function goToApiKeyPage() {
        this.setPage(0);
      },
      goToEnvsPage: function goToEnvsPage() {
        this.apiCall('/environments', function (environments) {
          if (environments.errors) {
            this.set('errors', environments.errors.map(function (err) {
              return err.message;
            }));
            return;
          }

          var envs = environments.data.filter(function (env) {
            return env.serviceConnection.type.toLowerCase() === 'cloudca';
          });

          this.environmentsByCodeName = envs.reduce(function (m, e) {
            m[e.serviceConnection.serviceCode + e.name] = e;
            return m;
          }, {});
          var envOpts = envs.map(function (env) {
            return {
              name: env.name,
              value: env.serviceConnection.serviceCode + env.name,
              group: env.serviceConnection.serviceCode
            };
          });
          this.set('environmentOptions', envOpts);

          if (this.config.serviceCode && this.config.environmentName) {
            this.set('config.environmentCodeName', this.config.serviceCode + this.config.environmentName);
          } else {
            this.set('config.environmentCodeName', envOpts[0].value);
          }
          this.setPage(1);
        }.bind(this));
      },
      goToComputePage: function goToComputePage() {
        var env = this.environmentsByCodeName[this.get('config.environmentCodeName')];
        if (env) {
          this.set('config.environmentName', env.name);
          this.set('config.serviceCode', env.serviceConnection.serviceCode);
          this.loadNetworks();
          this.loadTemplates();
          this.loadComputeOfferings();
          this.loadDiskOfferings();
        }
        var errs = this.get('errors');
        if (errs && errs.length > 0) {
          return;
        }
        this.setPage(2);
      },
      goToNodeTemplatePage: function goToNodeTemplatePage() {
        this.setPage(3);
      }
    },

    loadNetworks: function loadNetworks() {
      this.apiCall(this.getServicesApiEndpoint('networks'), function (listNetworksResponse) {
        if (listNetworksResponse.errors) {
          this.set('errors', listNetworksResponse.errors.map(function (err) {
            return err.message;
          }));
          return;
        }
        var networks = listNetworksResponse.data;
        if (!this.config.networkId && networks.length > 0) {
          this.set('config.networkId', networks[0].id);
        }
        this.set('networkOptions', networks.map(function (network) {
          return {
            name: network.name,
            value: network.id,
            group: network.vpcName
          };
        }));
      }.bind(this));
    },

    loadComputeOfferings: function loadComputeOfferings() {
      this.apiCall(this.getServicesApiEndpoint('computeofferings'), function (listComputeOfferingsResponse) {
        var _this = this;

        if (listComputeOfferingsResponse.errors) {
          this.set('errors', listComputeOfferingsResponse.errors.map(function (err) {
            return err.message;
          }));
          return;
        }
        var offerings = listComputeOfferingsResponse.data.filter(function (o) {
          return o.memoryInMB >= 2048 || o.custom;
        });
        this.set('computeOfferingOptions', offerings.map(function (o) {
          return {
            name: o.name,
            value: o.id,
            custom: o.custom
          };
        }));

        if (this.config.computeOffering) {
          var co = offerings.find(function (co) {
            return co.id === _this.config.computeOffering;
          });
          if (co && co.custom) {
            this.set("config.customCompute", true);
          }
        }

        if (!this.config.computeOffering && offerings.length > 0) {
          var off = offerings[0];
          this.set('config.computeOffering', off.id);
          if (off.custom) {
            this.set("config.customCompute", true);
            this.set("config.cpuCount", 1);
            this.set("config.memoryMb", 2048);
          }
        }
      }.bind(this));
    },

    loadDiskOfferings: function loadDiskOfferings() {
      this.apiCall(this.getServicesApiEndpoint('diskofferings'), function (listDiskOfferingsResponse) {
        var _this2 = this;

        if (listDiskOfferingsResponse.errors) {
          this.set('errors', listDiskOfferingsResponse.errors.map(function (err) {
            return err.message;
          }));
          return;
        }
        var offeringOptions = listDiskOfferingsResponse.data.map(function (offering) {
          return {
            name: offering.name,
            value: offering.id,
            customSize: offering.customSize,
            customIops: offering.customIops
          };
        });
        offeringOptions.push({
          name: "No additional disk",
          value: ""
        });

        if (this.config.additionalDiskOffering) {
          var co = offeringOptions.find(function (co) {
            return co.value === _this2.config.additionalDiskOffering;
          });
          if (co && co.customSize) {
            this.set("config.customDiskOffering", true);
          }
        }

        this.set('diskOfferingOptions', offeringOptions);
        if (!this.config.additionalDiskOffering && offeringOptions.length > 0) {
          var off = offeringOptions[0];
          this.set('config.additionalDiskOffering', off.value);
          this.set("config.customDiskOffering", off.customSize);
          if (off.customSize) {
            this.set("config.additionalDiskSizeGb", 20);
          }
        }
      }.bind(this));
    },

    loadTemplates: function loadTemplates() {
      this.apiCall(this.getServicesApiEndpoint('templates'), function (listTemplatesResponse) {
        if (listTemplatesResponse.errors) {
          this.set('errors', listTemplatesResponse.errors.map(function (err) {
            return err.message;
          }));
          return;
        }
        var removeTemplateRegex = /windows|centos 6|debian/i;
        var templates = listTemplatesResponse.data.filter(function (template) {
          return !template.name.match(removeTemplateRegex);
        });

        var templateOpts = templates.map(function (template) {
          return {
            name: template.name,
            value: template.id,
            group: template.isPublic ? 'Standard' : 'User defined',
            resizable: template.resizable,
            maxSizeInGb: template.maxSizeInGb,
            stepSizeInGb: template.stepSizeInGb,
            sshUser: template.defaultUsername
          };
        }).sortBy('group', 'name');
        if (!this.config.template && templateOpts.length > 0) {
          this.set('config.template', templateOpts[0].value);
        }

        this.set('templateOptions', templateOpts);
      }.bind(this));
    },

    apiCall: function apiCall(endpoint, callback) {
      var url = this.get('config.apiUrl') + endpoint;
      var xhr = new XMLHttpRequest();
      xhr.addEventListener('load', function () {
        callback(JSON.parse(this.responseText));
      });
      xhr.open('get', url, true);
      xhr.setRequestHeader('MC-Api-Key', this.get('config.apiKey'));
      xhr.setRequestHeader('Content-Type', 'application/json');
      xhr.send();
    },

    getServicesApiEndpoint: function getServicesApiEndpoint(entity) {
      return '/services/' + this.get('config.serviceCode') + '/' + this.get('config.environmentName') + '/' + entity;
    }

  });
});;
define('ui/components/driver-kt/component', ['exports', 'nodes/components/driver-kt/component'], function (exports, _component) {
  'use strict';

  Object.defineProperty(exports, "__esModule", {
    value: true
  });
  Object.defineProperty(exports, 'default', {
    enumerable: true,
    get: function () {
      return _component.default;
    }
  });
});in
