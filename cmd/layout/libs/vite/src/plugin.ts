import { join, resolve } from 'path';
import { PackageJson } from 'type-fest';
import { getWorkspacePackages, getWorkspaceRootPath } from './utils';

/**
 * Create alias for vite
 *
 * @param rootPaths - Root paths to looking for
 * @param exclude - Exclude apps, libs or packages by is package.json name
 */
export function workspacesAlias(rootPaths: string[], exclude: string[] = []) {
  return {
    name: 'vite-plugin-workspace-alias',

    config: (userConfig: Record<string, any>) => {
      const { alias = {} } = userConfig.resolve || {};

      const modifiedConfig = {
        ...userConfig,
        resolve: {
          alias: {
            ...Object.fromEntries(
              rootPaths
                .flatMap((rootPath) => getWorkspacePackages(rootPath))
                .filter(([_, pkg]) => {
                  const { name = null } = (pkg || {}) as PackageJson;

                  return name && !exclude.includes(name);
                })
                .map(([folder, pkg]) => {
                  const dir = folder as string;
                  const pkgJson = pkg as PackageJson & { source: string };

                  return [pkgJson.name, join(dir, pkgJson.source)];
                })
            ),
            ...alias
          }
        }
      };

      console.log('Automatic aliases:', modifiedConfig.resolve.alias);

      return modifiedConfig;
    }
  };
}

export function workspaceRollupOptions() {
  return {
    name: 'vite-plugin-workspace-rollup',
    config: (userConfig: Record<string, any>) => {
      const { build } = userConfig || {};
      const { rollupOptions = {} } = build || {};

      const packageRoot = resolve(getWorkspaceRootPath(), userConfig.root);
      const pkg = require(resolve(packageRoot, 'package.json')) as PackageJson;

      const { config = {} } = pkg;
      const { rollupOptions: rollup = {}, build: buildOptions = {} } = config;
      const { input = null } = rollup as any;

      if (input) {
        Object.keys(input).forEach((key) => {
          const indexEntry = input[key];
          input[key] = resolve(packageRoot, indexEntry);
        });
      }

      const settings = {
        ...userConfig,
        build: {
          ...build,
          ...(buildOptions as Record<string, any>),
          rollupOptions: {
            ...rollupOptions,
            ...(rollup as Record<string, any>)
          }
        }
      };

      return settings;
    }
  };
}
