'use strict';

const path = require('path');

const Descriptor = require('./descriptor').Descriptor;
const ServiceDescriptor = require('./service').Service; 
const deployment = exports;

const INSTALL_ROOT = path.resolve(path.join(
  __dirname,
  '../..'
));

class Deployment extends Descriptor {
  constructor (path, contents) {
    super('deployment', path, contents);
    this._descriptor = null;
  }

  /**
   * Returns a Descriptor instance for this deployment. Currently only returns
   * the default descriptor.
   *
   * Later on this would take a version number or something.
   */
  descriptor () {
    return new Promise((resolve, reject) => {

      if (this._descriptor) {
        return resolve(this._descriptor);
      }

      var descriptors = this.get('descriptors.supported');
      var defaultDescriptor = this.get('descriptors.default');

      var descriptor = descriptors.find((descriptor) => {
        return descriptor.version === defaultDescriptor; 
      });

      if (!descriptor) {
        return reject(new Error('Couldn\'t find default descriptor'));
      }

      var filePath = descriptor.url.replace('$ARIGATO_INSTALL', INSTALL_ROOT);
      filePath = filePath.replace('file://','');

      return ServiceDescriptor.read(filePath).then((descriptor) => {
        this._descriptor = descriptor;
        resolve(descriptor);
      }).catch(reject);
    });
  }

  static read (filePath) {
    return super.read(Deployment, 'deployment', filePath);
  }
}

deployment.Deployment = Deployment;