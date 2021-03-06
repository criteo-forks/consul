export default {
  Name: (item, value) => item.Name,
  Description: (item, value) => item.Description,
  AccessorID: (item, value) => item.AccessorID,
  Role: (item, value) => (item.Roles || []).map(item => item.Name),
  Policy: (item, value) => {
    return (item.Policies || [])
      .map(item => item.Name)
      .concat((item.ServiceIdentities || []).map(item => item.ServiceName))
      .concat((item.NodeIdentities || []).map(item => item.NodeName));
  },
};
