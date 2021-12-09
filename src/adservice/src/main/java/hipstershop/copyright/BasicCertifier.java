package hipstershop.copyright;

import hipstershop.Demo;

import java.util.ArrayList;
import java.util.List;

class BasicCertifier implements Certifier {
    @Override
    public List<Demo.Ad> certify(List<Demo.Ad> ads) {
        List<Demo.Ad> result = new ArrayList<>(ads.size());
        for (Demo.Ad ad : ads) {
            if (!ad.getText().contains("©")) {
                result.add(ad);
            }
        }
        return result;
    }
}
