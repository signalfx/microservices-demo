package hipstershop.copyright;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CopyrightCertification {

    private final Map<String, Certifier> certifiers = new HashMap<>();

    {
        certifiers.put("vintage", new BasicCertifier());
        certifiers.put("cycling", new BasicCertifier());
        certifiers.put("cookware", new BasicCertifier());
        certifiers.put("photography", new PhotographyCertifier());
        certifiers.put("gardening", new BasicCertifier());
    }

    public List<hipstershop.Demo.Ad> certify(String category, List<hipstershop.Demo.Ad> ads) {
        Certifier certifier = certifiers.get(category);
        return certifier.certify(ads);
    }


}
