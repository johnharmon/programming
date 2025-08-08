# lookup_plugins/macgen.py
from __future__ import absolute_import, division, print_function
__metaclass__ = type

import random

from ansible.plugins.lookup import LookupBase
from ansible.errors import AnsibleError


def random_mac():
    """Generate a locally administered unicast MAC address."""
    mac = [0x02, random.randint(0x00, 0x7F), random.randint(0x00, 0xFF),
           random.randint(0x00, 0xFF), random.randint(0x00, 0xFF), random.randint(0x00, 0xFF)]
    return ":".join("{:02x}".format(octet) for octet in mac)


class LookupModule(LookupBase):

    def run(self, terms, variables=None, **kwargs): # type: ignore
        count = kwargs.get('count')
        exclude = set(kwargs.get('exclude', []))

        if count is None:
            raise AnsibleError("The 'count' parameter is required.")

        result = []
        attempts = 0
        max_attempts = 1000

        while len(result) < count and attempts < max_attempts:
            mac = random_mac()
            if mac not in exclude and mac not in result:
                result.append(mac)
            attempts += 1

        if len(result) < count:
            raise AnsibleError(f"Could only generate {len(result)} unique MACs after {max_attempts} attempts.")

        return result
