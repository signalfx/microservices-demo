package hipstershop.copyright;

import hipstershop.Demo;

import java.util.List;

// Filters a list of ads to only those that pass copyright certification
interface Certifier {
    List<Demo.Ad> certify(List<Demo.Ad> ads);
}
