// Create a function that get the content after explore and store it in an array, like this : '/explore/all' => ['all']
export function getExploreContent(url: string) {
  return url.split('/explore').slice(-1);
}

export enum LinksDefinition {
  all = 'Explore all',
  user = 'User',
}
